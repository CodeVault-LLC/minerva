import z from "zod";

export const obfuscations = z.enum([
  "eval",
  "function",
  "window",
  "hex",
  "unicode",
  "hexadecimal",
  "anonymous",
]);

export const obfuscationPatterns = [
  {
    name: obfuscations.Enum.eval,
    patterns: [/eval\(/g],
  },
  {
    name: obfuscations.Enum.function,
    patterns: [/Function\(/g],
  },
  {
    name: obfuscations.Enum.window,
    patterns: [/window\["[^"]+"\]/g],
  },
  {
    name: obfuscations.Enum.hex,
    patterns: [/\\x[0-9A-Fa-f]{2}/g],
  },
  {
    name: obfuscations.Enum.unicode,
    patterns: [/\\u[0-9A-Fa-f]{4}/g],
  },
  {
    name: obfuscations.Enum.hexadecimal,
    patterns: [/_0x[0-9a-f]{4,}/g],
  },
  {
    name: obfuscations.Enum.anonymous,
    patterns: [/function\s*\(\s*[^)]*\)\s*\{[^}]*\}/g],
  },
];
