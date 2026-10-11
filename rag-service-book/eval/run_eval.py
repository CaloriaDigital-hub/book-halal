"""
Оценка поиска RAG по набору вопросов с известными страницами ответа.

Ходит в POST /debug/retrieve/{book_id} (только поиск, без LLM) и считает,
на каком месте в выдаче оказалась нужная страница.

    python eval/run_eval.py --book <book_id> --label baseline
    python eval/run_eval.py --book <book_id> --label chunks --baseline eval/results/baseline.json
"""
import argparse
import json
import pathlib
import statistics
import time
import urllib.request

HERE = pathlib.Path(__file__).parent
KS = (1, 3, 5, 10)


def retrieve(url: str, book_id: str, question: str, top_k: int) -> list[dict]:
    req = urllib.request.Request(
        f"{url}/debug/retrieve/{book_id}",
        data=json.dumps({"question": question, "top_k": top_k}).encode(),
        headers={"Content-Type": "application/json"},
    )
    with urllib.request.urlopen(req, timeout=120) as res:
        return json.load(res)["results"]


def first_hit(results: list[dict], pages: list[int]) -> int | None:
    """Позиция (с 1) первого результата с нужной страницы."""
    for i, r in enumerate(results, 1):
        if r["page"] in pages:
            return i
    return None


def score(rows: list[dict]) -> dict[str, float]:
    n = len(rows)
    metrics = {f"hit@{k}": sum(1 for r in rows if r["rank"] and r["rank"] <= k) / n for k in KS}
    metrics["mrr"] = sum(1 / r["rank"] for r in rows if r["rank"]) / n
    return metrics


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--url", default="http://localhost:8000")
    ap.add_argument("--book", required=True)
    ap.add_argument("--label", required=True, help="имя прогона, результат ляжет в eval/results/<label>.json")
    ap.add_argument("--top-k", type=int, default=10)
    ap.add_argument("--baseline", help="json прошлого прогона для сравнения по вопросам")
    args = ap.parse_args()

    questions = json.loads((HERE / "questions.json").read_text(encoding="utf-8"))
    rows = []
    t0 = time.time()
    for q in questions:
        results = retrieve(args.url, args.book, q["question"], args.top_k)
        rank = first_hit(results, q["pages"])
        rows.append({**q, "rank": rank, "got": [r["page"] for r in results]})
        print(f"{q['id']:>3} {'—' if rank is None else rank:>2}  {q['question']}")

    metrics = score(rows)
    metrics["seconds_per_query"] = (time.time() - t0) / len(rows)
    for name in sorted({r.get("set", "base") for r in rows}):
        for k, v in score([r for r in rows if r.get("set", "base") == name]).items():
            metrics[f"{name}/{k}"] = v

    print()
    for k, v in metrics.items():
        print(f"  {k:<22} {v:.2f}")

    out = HERE / "results" / f"{args.label}.json"
    out.parent.mkdir(exist_ok=True)
    out.write_text(json.dumps({"label": args.label, "metrics": metrics, "rows": rows}, ensure_ascii=False, indent=1), encoding="utf-8")
    print(f"сохранено: {out}")

    if args.baseline:
        base = json.loads(pathlib.Path(args.baseline).read_text(encoding="utf-8"))
        print(f"\nсравнение с {base['label']}:")
        for k, v in metrics.items():
            b = base["metrics"].get(k)
            if b is not None:
                print(f"  {k:<22} {b:.2f} → {v:.2f} ({v - b:+.2f})")
        base_rank = {r["id"]: r["rank"] for r in base["rows"]}
        changed = [(r, base_rank.get(r["id"])) for r in rows if r["rank"] != base_rank.get(r["id"])]
        for r, before in changed:
            fmt = lambda x: "—" if x is None else str(x)
            print(f"  #{r['id']:<3} {fmt(before):>2} → {fmt(r['rank']):<2} {r['question']}")


if __name__ == "__main__":
    main()
