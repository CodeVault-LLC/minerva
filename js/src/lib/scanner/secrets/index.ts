import type { Secret } from "~/lib/types";
import { SECRET_PATTERNS } from "~/lib/scanner/secrets/lists";

export const scanSecrets = (
  scripts: {
    content: string;
    src: string;
  }[]
): Secret[] => {
  const foundSecrets: Secret[] = [];

  scripts.forEach((script) => {
    SECRET_PATTERNS.forEach((pattern) => {
      const matches = regexHandler(script, pattern.pattern);
      if (matches.length > 0) {
        foundSecrets.push({
          name: pattern.name,
          results: matches,
        });
      }
    });
  });

  return foundSecrets;
};

export interface RegexResult {
  match: string;
  line: number;
  script: string;
}

export const regexHandler = (
  script: { content: string; src: string },
  regex: RegExp
): RegexResult[] | never[] => {
  const message = script.content;
  const lines = message.split("\n");
  const results: RegexResult[] = [];

  for (let lineNumber = 0; lineNumber < lines.length; lineNumber++) {
    const line = lines[lineNumber];
    const matches = line.matchAll(regex);

    for (const match of matches) {
      results.push({
        match: match[0],
        line: lineNumber + 1, // Line numbers are usually 1-based
        script: script.src,
      });
    }
  }

  return results.length ? results : [];
};
