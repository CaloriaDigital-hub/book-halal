import pickle
import os
import numpy as np
import networkx as nx
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

model = SentenceTransformer("sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2")

GRAPH_DIR = "graphs"
os.makedirs(GRAPH_DIR, exist_ok=True)


def _graph_path(book_id: str) -> str:
    return os.path.join(GRAPH_DIR, f"{book_id}.pkl")


def build_graph(book_id: str, sentences: list[str], pages: list[int], threshold: float = 0.65):
    """
    Строим граф где узел = предложение, ребро = семантическая близость.
    """
    embeddings = model.encode(sentences)
    similarity_matrix = cosine_similarity(embeddings)

    G = nx.Graph()

    # Добавляем узлы
    for i, (sentence, page) in enumerate(zip(sentences, pages)):
        G.add_node(i, text=sentence, page=page, embedding=embeddings[i])

    # Добавляем рёбра по порогу сходства
    for i in range(len(sentences)):
        for j in range(i + 1, len(sentences)):
            score = float(similarity_matrix[i][j])
            if score >= threshold:
                G.add_edge(i, j, weight=score)

    # Сохраняем граф локально
    with open(_graph_path(book_id), "wb") as f:
        pickle.dump(G, f)

    return G


def load_graph(book_id: str) -> nx.Graph | None:
    path = _graph_path(book_id)
    if not os.path.exists(path):
        return None
    with open(path, "rb") as f:
        return pickle.load(f)


def search_graph(book_id: str, query: str, top_k: int = 5, query_threshold: float = 0.4) -> list[dict]:
    G = load_graph(book_id)
    if G is None or len(G.nodes) == 0:
        return []

    query_embedding = model.encode([query])[0]

    scores = {}
    for node_id, data in G.nodes(data=True):
        score = float(cosine_similarity([query_embedding], [data["embedding"]])[0][0])
        scores[node_id] = score

    sorted_scores = sorted(scores.items(), key=lambda x: x[1], reverse=True)
    best_node = sorted_scores[0][0]

    # Соседей берём только если они сами релевантны вопросу
    context_nodes = {best_node}
    for neighbor in G.neighbors(best_node):
        edge_weight = G[best_node][neighbor]["weight"]
        neighbor_query_score = scores[neighbor]
        if edge_weight >= 0.65 and neighbor_query_score >= query_threshold:
            context_nodes.add(neighbor)

    results = []
    for node_id in context_nodes:
        data = G.nodes[node_id]
        results.append({
            "text": data["text"],
            "page": data["page"],
            "score": scores[node_id]
        })

    results.sort(key=lambda x: x["score"], reverse=True)
    return results[:top_k]

def delete_graph(book_id: str):
    path = _graph_path(book_id)
    if os.path.exists(path):
        os.remove(path)