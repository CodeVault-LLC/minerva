import type { Subscription } from "@prisma/client";
import { prisma } from "~/data-source";

interface CreateSubscriptionProps {
  userId: number;
  plan: string;
  stripeSessionId: string;
  price: number;
  status: string;
}

export const createSubscription = async ({
  userId,
  plan,
  stripeSessionId,
  price,
  status,
}: CreateSubscriptionProps): Promise<Subscription> => {
  const subscription = await prisma.subscription.create({
    data: {
      userId,
      price,
      status,
      subscriptionId: stripeSessionId,
      type: plan,
    },
  });

  return subscription;
};

export const updateSubscription = async (
  subscriptionId: string,
  status: string
): Promise<Subscription> => {
  const subscription = await prisma.subscription.update({
    where: {
      subscriptionId,
    },
    data: {
      status,
    },
  });

  return subscription;
};

export const getUserSubscription = async (
  userId: number
): Promise<Subscription | null> => {
  const subscription = await prisma.subscription.findFirst({
    where: {
      userId,
    },
  });

  return subscription;
};
