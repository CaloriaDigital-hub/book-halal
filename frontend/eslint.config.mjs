import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";

// Границы модульного монолита:
//   app      → может брать modules/* (только публичный index / server) и shared
//   modules  → только shared и публичные API других модулей, внутри модуля — относительные импорты
//   shared   → ничего не знает ни про modules, ни про app
const DEEP_MODULE_IMPORT = {
  regex: "^@/modules/[^/]+/(?!server$).+",
  message: "Импортируй модуль через его публичный API: @/modules/<name> или @/modules/<name>/server.",
};
const RELATIVE_ESCAPE = {
  regex: "^\\.\\./\\.\\./\\.\\./",
  message: "Не выходи за пределы модуля относительным путём — используй публичный API через @/.",
};

const eslintConfig = defineConfig([
  ...nextVitals,
  ...nextTs,
  {
    files: ["src/app/**/*.{ts,tsx}", "src/proxy.ts"],
    rules: { "no-restricted-imports": ["error", { patterns: [DEEP_MODULE_IMPORT] }] },
  },
  {
    files: ["src/modules/**/*.{ts,tsx}"],
    rules: {
      "no-restricted-imports": [
        "error",
        {
          patterns: [
            DEEP_MODULE_IMPORT,
            RELATIVE_ESCAPE,
            { regex: "^@/app(/|$)", message: "Модули не зависят от app — роуты передаются пропсами." },
          ],
        },
      ],
    },
  },
  {
    files: ["src/shared/**/*.{ts,tsx}"],
    rules: {
      "no-restricted-imports": [
        "error",
        { patterns: [{ regex: "^@/(modules|app)(/|$)", message: "shared не зависит от modules и app." }] },
      ],
    },
  },
  // Override default ignores of eslint-config-next.
  globalIgnores([
    // Default ignores of eslint-config-next:
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
  ]),
]);

export default eslintConfig;
