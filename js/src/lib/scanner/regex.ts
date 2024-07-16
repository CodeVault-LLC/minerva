// Every secret should be looking like this:
/*
  [
  {
    script: "https://example.com/script.js",
    checks: [
      {
        name: "Secrets",
        matches: [{
          name: "AWS Key",
          line: 1,
          match: "AKIA...",
        }]
      },
      {
        name: "Databases",
        matches: [{
          name: "Firebase",
          line: 1,
          match: "firebase",
          }]
      ]

    }
    ]
  */

import { scanDatabase } from "./database";
import { scanObfuscations } from "./obfuscation";
import { scanSecrets } from "./secrets";

export interface RegexResponse {
  name: string;
  line: number;
  match: string;
}

export interface Regexer {
  script: string;
  checks: {
    name: string;
    matches: {
      name: string;
      line: number;
      match: string;
    }[];
  }[];
}

export const regexHandler = (
  scripts: {
    src: string;
    content: string;
  }[]
): Regexer[] => {
  const scans = [
    {
      name: "Findings",
      function: scanSecrets,
    },
    {
      name: "Databases",
      function: scanDatabase,
    },
    {
      name: "Obfuscations",
      function: scanObfuscations,
    },
  ];

  const results: Regexer[] = [];

  scripts.forEach((script) => {
    const checks: {
      name: string;
      matches: RegexResponse[];
    }[] = [];

    scans.forEach((scan) => {
      const found = scan.function([script]);
      if (found.length > 0) {
        checks.push({
          name: scan.name,
          matches: found[0].results.map((result) => ({
            name: result.match,
            line: result.line,
            match: result.match,
          })),
        });
      }
    });

    results.push({
      script: script.src,
      checks,
    });
  });
};
