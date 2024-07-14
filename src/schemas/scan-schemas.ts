import { z } from "zod";

export const scanSchema = z.object({
  url: z.string(),
  depth: z.number().optional().default(2), // from 1 to 5
  doScripts: z.boolean().optional().default(true),
  doStyles: z.boolean().optional().default(false),
  doImages: z.boolean().optional().default(false),
  doLinks: z.boolean().optional().default(true),
  scripts: z
    .array(
      z.object({
        src: z.string(),
        content: z.string(),
      })
    )
    .optional()
    .default([]),
});
