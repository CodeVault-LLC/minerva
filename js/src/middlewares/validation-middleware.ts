import { type Request, type Response, type NextFunction } from "express";
import { type z, ZodError } from "zod";
import { StatusCodes } from "http-status-codes";

export const validateData = (schema: z.ZodSchema) => {
  return (req: Request, res: Response, next: NextFunction) => {
    try {
      schema.parse(req.body);
      next();
    } catch (e) {
      if (e instanceof ZodError) {
        return res.status(StatusCodes.BAD_REQUEST).json({
          status: "error",
          errors: e.errors,
        });
      }
      next(e);
    }
  };
};
