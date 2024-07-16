import { regexHandler, type RegexResult } from "../secrets";
import { obfuscationPatterns, obfuscations } from "./lists";

export const scanObfuscations = (
  scripts: {
    src: string;
    content: string;
  }[]
) => {
  const foundObfuscations: {
    name: string;
    results: RegexResult[];
  }[] = [];

  scripts.forEach((script) => {
    obfuscationPatterns.forEach((regex) => {
      regex.patterns.forEach((pattern) => {
        const matches = regexHandler(script, pattern);
        if (matches.length > 0) {
          foundObfuscations.push({
            name: regex.name,
            results: matches,
          });

          try {
            obfuscations.parse(regex.name);
          } catch (error) {
            // eslint-disable-next-line no-console -- Allowed for logging
            console.error(error);
          }
        }
      });
    });
  });

  return foundObfuscations;
};
