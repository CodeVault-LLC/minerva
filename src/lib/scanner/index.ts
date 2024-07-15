import { scanSchema } from "~/schemas/scan-schemas";
import type { Secret } from "../types";
import { type RegexResult, scanSecrets } from "./secrets";
import { scanDatabase, scanDatabaseLogins } from "./database";
import { scanObfuscations } from "./obfuscation";

interface Output {
  secrets: Secret[];
  // client sided databases such as Firebase, AWS, etc
  databases: {
    name: string;
    results: RegexResult[];
  }[];
  databaseLogins: {
    name: string;
    results: RegexResult[];
  }[];
  // obfuscation methods
  obfuscations: {
    name: string;
    results: RegexResult[];
  }[];
}

export const analyze = async (body): Promise<Output> => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars -- Ignore unused variables
  const { url, depth, doScripts, doStyles, doImages, doLinks, scripts } =
    scanSchema.parse(body);

  const foundSecrets = scanSecrets(scripts);
  const databases = scanDatabase(scripts);
  const obfuscations = scanObfuscations(scripts);
  const databaseLogins = scanDatabaseLogins(scripts);

  return {
    secrets: foundSecrets,
    databases,
    databaseLogins,
    obfuscations,
  };

  // known bad scripts and urls check
  // known bad domains check
  // known database check
  // known library checker
  // known bad file extensions
};
