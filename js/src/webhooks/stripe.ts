import { Router } from "express";
import Stripe from "stripe";
import bodyParser from "body-parser";
import {
  createSubscription,
  updateSubscription,
} from "~/controller/subscription.controller";
import { getUserByEmail } from "~/controller/user.controller";

const stripe = new Stripe(process.env.STRIPE_SECRET_KEY ?? "");

const endpointSecret = process.env.STRIPE_WEBHOOK_SECRET ?? "";

const app = Router();

app.post(
  "/",
  bodyParser.raw({ type: "application/json" }),
  // eslint-disable-next-line @typescript-eslint/no-misused-promises -- This is a callback function
  async (req, res) => {
    const sig = req.headers["stripe-signature"];

    let event: Stripe.Event;

    try {
      event = stripe.webhooks.constructEvent(
        // eslint-disable-next-line @typescript-eslint/no-unsafe-argument -- This is a callback function
        req.body,
        sig ?? "",
        endpointSecret
      );
    } catch (err) {
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions, @typescript-eslint/no-unsafe-member-access -- This is an error message
      return res.status(400).send(`Webhook Error: ${err.message}`);
    }

    // Handle the event
    if (event.type === "customer.subscription.created") {
      const subscription = event.data.object;
      const customerId = subscription.customer as string;

      try {
        const user = (await stripe.customers.retrieve(
          customerId
        )) as Stripe.Customer;

        const dbUser = await getUserByEmail(user.email ?? "");

        if (!dbUser) {
          return res.status(404).json({ received: false });
        }

        if (dbUser.subscriptions.length > 0) {
          // this user already has a subscription
          return res.status(400).json({ received: false });
        }

        const productName = (await stripe.products.retrieve(
          subscription.items.data[0].price.product as string
        )) as Stripe.Product;

        await createSubscription({
          userId: dbUser.id,
          plan: productName.name,
          stripeSessionId: subscription.id,
          price: subscription.items.data[0].price.unit_amount ?? 0,
          status: subscription.status,
        });
        return res.json({ received: true });
      } catch (err) {
        return res.status(400).json({ received: false });
      }
    } else if (event.type === "customer.subscription.updated") {
      const subscription = event.data.object;
      const customerId = subscription.customer as string;

      try {
        const user = (await stripe.customers.retrieve(
          customerId
        )) as Stripe.Customer;

        const dbUser = await getUserByEmail(user.email ?? "");

        if (!dbUser) {
          return res.status(404).json({ received: false });
        }

        await updateSubscription(subscription.id, subscription.status);
        return res.json({ received: true });
      } catch (err) {
        return res.status(400).json({ received: false });
      }
    } else if (event.type === "customer.subscription.deleted") {
      try {
        const subscription = event.data.object;

        await updateSubscription(subscription.id, subscription.status);
        return res.json({ received: true });
      } catch (err) {
        return res.status(400).json({ received: false });
      }
    } else {
      res.json({ received: true });
    }
  }
);

export { app as stripeWebhook };
