import { regexHandler, type RegexResult } from "../secrets";
import { ClientDatabases, regexes } from "./lists";

export const scanDatabase = (
  scripts: {
    src: string;
    content: string;
  }[]
): {
  name: string;
  results: RegexResult[];
}[] => {
  // Look for certain patterns in the script

  // for each script go through the contents and then run the regex checks.
  const foundDatabases: {
    name: string;
    results: RegexResult[];
  }[] = [];

  scripts.forEach((script) => {
    regexes.forEach((regex) => {
      regex.patterns.forEach((pattern) => {
        const matches = regexHandler(script, pattern);
        if (matches.length > 0) {
          foundDatabases.push({
            name: regex.name,
            results: matches,
          });

          // verify the schema
          try {
            ClientDatabases.parse(regex.name);
          } catch (error) {
            // eslint-disable-next-line no-console -- Allowed for logging
            console.error(error);
          }
        }
      });
    });
  });

  return foundDatabases;
};
