import { json, Router } from "express";
import { analyze } from "~/lib/scanner";
import { validateData } from "~/middlewares/validation-middleware";
import { scanSchema } from "~/schemas/scan-schemas";

const router = Router();

router.post(
  "/",
  json({ limit: "50mb" }),
  validateData(scanSchema),
  // eslint-disable-next-line @typescript-eslint/no-misused-promises -- This is a callback function
  async (req, res) => {
    const { url, depth, doScripts, doStyles, doImages, doLinks, scripts } =
      scanSchema.parse(req.body);

    // Analyze everything.
    const analysis = await analyze({
      url,
      depth,
      doScripts,
      doStyles,
      doImages,
      doLinks,
      scripts,
    });

    res.json(analysis);
  }
);

export const scanRouter = router;
