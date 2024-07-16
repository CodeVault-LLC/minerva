import { PrismaClient } from "@prisma/client";
import consolaLogger from "consola";

const prisma = new PrismaClient();

prisma
  .$connect()
  .then(() => {
    consolaLogger.success("Prisma connected successfully");
  })
  .catch((error: unknown) => {
    consolaLogger.error("Prisma connection error:", error);
    process.exit(1);
  });

export { prisma };
