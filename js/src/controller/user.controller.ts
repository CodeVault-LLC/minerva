import type { Subscription, User } from "@prisma/client";
import { prisma } from "~/data-source";

export const getUsers = async (): Promise<User[]> => {
  return prisma.user.findMany();
};

export const getUserById = async (id: number): Promise<User | null> => {
  return prisma.user.findUnique({
    where: {
      id,
    },
  });
};

export const createUser = async (user: {
  discordId: string;
  username: string;
  email: string;
  accessToken: string;
  avatar: string;
  createdAt: Date;
  provider: string;
  updatedAt: Date;
}): Promise<User> => {
  return prisma.user.create({
    data: user,
  });
};

export const findUserByDiscordId = async (
  discordId: string
): Promise<User | null> => {
  return prisma.user.findFirst({
    where: {
      discordId,
    },
  });
};

export const getUserByEmail = async (
  email: string
): Promise<(User & { subscriptions: Subscription[] }) | null> => {
  return prisma.user.findFirst({
    where: {
      email,
    },
    include: {
      subscriptions: true,
    },
  });
};
