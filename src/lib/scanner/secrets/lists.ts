import { Category, SubCategory } from "~/lib/constants";
import type { RegexPattern } from "~/lib/types";

export const SECRET_PATTERNS: RegexPattern[] = [
  {
    pattern:
      /(?:key|token|secret|password|apiKey|authToken|accessToken|sessionKey|privateKey|publicKey)[\s:="']+(?:[\w-]+)/gi,
    name: "Sensitive Data Keywords",
    description:
      "Matches common keywords related to sensitive data like keys, tokens, secrets, etc.",
    category: Category.SENSITIVE_DATA,
    subCategory: SubCategory.TOKEN,
  },
  {
    pattern: /(?:"|')\w{20,}(?:"|')/gi,
    name: "Long String Tokens",
    description: "Matches long strings that may resemble API keys or tokens.",
    category: Category.SENSITIVE_DATA,
    subCategory: SubCategory.API_KEY,
  },
  {
    pattern:
      /\b(?:secret|password|apiKey|token|access_token|auth_token|credentials)\b\s*[:=]\s*(?:"|')[^"']+/gi,
    name: "Sensitive Data Assignment",
    description:
      "Matches assignments of sensitive data like secrets, passwords, tokens, etc.",
    category: Category.SENSITIVE_DATA,
    subCategory: SubCategory.PASSWORD,
  },

  // Game Variables (Level, Health, Lives, etc.)
  {
    pattern:
      /\b(?:level|health|lives|damage|score|points|enemy|player|play|experience|mana|stamina|currency|coins|gems|ammo|bullets)\b\s*(?:=|[+-]=)\s*[^;]+/gi,
    name: "Game Variable Assignment",
    description:
      "Matches assignment or modification of common game-related variables.",
    category: Category.GAME_VARIABLE,
    subCategory: SubCategory.GAME_STAT,
  },
  {
    pattern:
      /\b(?:level|health|lives|damage|score|points|enemy|player|play|experience|mana|stamina|currency|coins|gems|ammo|bullets)\b/gi,
    name: "Game Variables",
    description: "Identifies common game-related variables in code.",
    category: Category.GAME_VARIABLE,
    subCategory: SubCategory.GAME_STAT,
  },

  // Function Names (Potentially Interesting or Vulnerable)
  {
    pattern:
      /\b(?:init|start|stop|pause|resume|render|update|draw|reset|destroy|remove|add|insert|delete|attack|move|jump|run|shoot|spawn|load|save)\b\s*\([^)]*\)/gi,
    name: "Game Function Calls",
    description: "Matches calls to common game-related functions.",
    category: Category.FUNCTION_NAME,
    subCategory: SubCategory.GAME_ACTION,
  },
  {
    pattern:
      /\b(?:render|update|draw|init|play|pause|resume|restart|stop|load|save|reset|destroy|remove|add|insert|delete|attack|move|jump|run|shoot|spawn)\b/gi,
    name: "Game Function Names",
    description: "Identifies common game-related function names in code.",
    category: Category.FUNCTION_NAME,
    subCategory: SubCategory.GAME_ACTION,
  },

  // Weird or Generic Names
  {
    pattern:
      /\b(?:data|result|value|temp|val|tmp|obj|item|element|component|list|array|dict|config|info|params|args|settings|state|counter|index|i|j|k|n|x|y|z)\b/gi,
    name: "Generic Names",
    description:
      "Matches common generic variable names which might be placeholders or less secure.",
    category: Category.GENERIC_NAME,
    subCategory: SubCategory.COMMON_VARIABLE,
  },

  // URLs, Endpoints, and Paths
  {
    pattern: /https?:\/\/[^\s/$.?#].[^\s"]*/gi,
    name: "URL Matching",
    description: "Matches HTTP and HTTPS URLs.",
    category: Category.URL_PATH,
    subCategory: SubCategory.URL,
  },
  {
    pattern: /(?:[a-zA-Z]:)?(?:\\|\/)[\w\s.\\/-]*\.[a-zA-Z]{2,4}/gi,
    name: "File Path Matching",
    description: "Matches file paths that can be either absolute or relative.",
    category: Category.URL_PATH,
    subCategory: SubCategory.FILE_PATH,
  },
  {
    pattern: /\b(?:endpoint|url|uri|path)\b\s*[:=]\s*(?:"|')[^"']+/gi,
    name: "Endpoint Assignment",
    description: "Matches assignments of URLs or paths to variables.",
    category: Category.URL_PATH,
    subCategory: SubCategory.URL,
  },

  // Potential Security Risks (eval, setTimeout, setInterval)
  {
    pattern: /\beval\s*\(.*?\)/gi,
    name: "Eval Usage",
    description:
      "Matches usage of eval function, which can be a security risk.",
    category: Category.SECURITY_RISK,
    subCategory: SubCategory.DOM_MANIPULATION,
  },
  {
    pattern: /\bsetTimeout\s*\(.*?\)/gi,
    name: "setTimeout Usage",
    description: "Matches usage of setTimeout function.",
    category: Category.SECURITY_RISK,
    subCategory: SubCategory.COMMON_VARIABLE,
  },
  {
    pattern: /\bsetInterval\s*\(.*?\)/gi,
    name: "setInterval Usage",
    description: "Matches usage of setInterval function.",
    category: Category.SECURITY_RISK,
    subCategory: SubCategory.COMMON_VARIABLE,
  },
  {
    pattern: /\binnerHTML\b\s*=\s*.*?;/gi,
    name: "Direct innerHTML Assignment",
    description: "Matches direct assignment to innerHTML property.",
    category: Category.SECURITY_RISK,
    subCategory: SubCategory.DOM_MANIPULATION,
  },
  {
    pattern: /\b(?:location|window\.location)\b\s*=\s*.*?;/gi,
    name: "Location Manipulation",
    description: "Matches assignment to location object.",
    category: Category.SECURITY_RISK,
    subCategory: SubCategory.DOM_MANIPULATION,
  },

  // Variable Declarations and Initializations
  {
    pattern: /\b(?:var|let|const)\s+\w+\s*=\s*[^;]+;/gi,
    name: "Variable Declaration",
    description: "Matches general variable declarations with initialization.",
    category: Category.VARIABLE_DECLARATION,
    subCategory: SubCategory.COMMON_VARIABLE,
  },
  {
    pattern: /\b(?:this|window)\.\w+\s*=\s*[^;]+;/gi,
    name: "Global or Window Variable Assignment",
    description: "Matches assignments to global or window-scoped variables.",
    category: Category.VARIABLE_DECLARATION,
    subCategory: SubCategory.COMMON_VARIABLE,
  },

  // JSON and Object Properties
  {
    pattern:
      /\b(?:property|key|value|json|config|data|metadata)\b\s*[:=]\s*(?:"|')[^"']+/gi,
    name: "JSON or Object Property",
    description: "Matches assignments to common JSON or object properties.",
    category: Category.JSON_PROPERTY,
    subCategory: SubCategory.COMMON_VARIABLE,
  },

  // Common JavaScript Patterns
  {
    pattern: /\b(?:console\.(?:log|warn|error|debug))\b\s*\(.*?\)/gi,
    name: "Console Logging",
    description: "Matches console logging functions.",
    category: Category.JS_PATTERN,
    subCategory: SubCategory.CONSOLE_LOG,
  },
  {
    pattern: /\b(?:for|while)\s*\([^)]*\)\s*\{[^}]*\}/gi,
    name: "Loops",
    description: "Matches for and while loops.",
    category: Category.JS_PATTERN,
    subCategory: SubCategory.LOOP,
  },
  {
    pattern:
      /\b(?:document\.|innerHTML|outerHTML|appendChild|insertAdjacentHTML|innerText|outerText|querySelector|querySelectorAll|getElementById|getElementsByClassName|getElementsByTagName|createElement)\b/gi,
    name: "DOM Manipulation",
    description: "Matches common DOM manipulation methods.",
    category: Category.JS_PATTERN,
    subCategory: SubCategory.DOM_MANIPULATION,
  },
  {
    pattern: /\bfetch\s*\(.*?\)/gi,
    name: "Fetch API Usage",
    description: "Matches usage of Fetch API.",
    category: Category.JS_PATTERN,
    subCategory: SubCategory.AJAX,
  },
  {
    pattern: /\b(?:axios|request|ajax|$.ajax|$.get|$.post)\s*\(.*?\)/gi,
    name: "AJAX Libraries",
    description: "Matches usage of common AJAX libraries and methods.",
    category: Category.JS_PATTERN,
    subCategory: SubCategory.AJAX,
  },

  // Comments (Potentially Sensitive Information)
  {
    pattern: /\/\*[^*]*\*+(?:[^/*][^*]*\*+)*\//gi,
    name: "Block Comments",
    description: "Matches block comments in the code.",
    category: Category.COMMENT,
    subCategory: SubCategory.BLOCK_COMMENT,
  },
  {
    pattern: /\/\/.*/gi,
    name: "Line Comments",
    description: "Matches single-line comments in the code.",
    category: Category.COMMENT,
    subCategory: SubCategory.LINE_COMMENT,
  },

  // Regular Expressions
  {
    pattern: /\/[^/]+\/[gimuy]*/gi,
    name: "Regex Literals",
    description: "Matches regular expression literals.",
    category: Category.REGEX_LITERAL,
  },
];
