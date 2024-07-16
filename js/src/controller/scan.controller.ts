import { prisma } from "~/data-source";
import { Secret } from "~/types/secret";

export const addSecret = async (secret: Secret) => {
  const { pattern } = secret;
}
