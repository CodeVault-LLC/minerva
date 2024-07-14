import express, { type Express } from "express";
import passport from "passport";
import session from "express-session";
import "./strategy/discord";
import Stripe from "stripe";
import { config } from "dotenv";
import type { User } from "@prisma/client";
import {
  createSubscription,
  getUserSubscription,
} from "./controller/subscription.controller";
import { stripeWebhook } from "./webhooks/stripe";
import { scanRouter } from "./routes/scan";

config();
const stripe = new Stripe(process.env.STRIPE_SECRET_KEY ?? "");

const app: Express = express();

app.use(
  session({
    secret: "secret-key",
    resave: false,
    saveUninitialized: false,
  })
);

app.use(passport.initialize());
app.use(passport.session());

// eslint-disable-next-line @typescript-eslint/no-unsafe-argument -- This is a callback function
app.get("/auth/discord", passport.authenticate("discord"));

app.get(
  "/auth/discord/callback",
  // eslint-disable-next-line @typescript-eslint/no-unsafe-argument -- This is a callback function
  passport.authenticate("discord", {
    failureRedirect: "/",
  }),
  (req, res) => {
    res.send("<script>window.close();</script>");
  }
);

// eslint-disable-next-line @typescript-eslint/no-misused-promises -- This is a callback function
app.get("/me", async (req, res) => {
  // if not user then send 401 if user then send 200
  if (!req.user) {
    res.status(401).send("Unauthorized");
    return;
  }

  // get their subscription status and plan.
  const subscription = await getUserSubscription((req.user as User).id);

  res.status(200).json({
    id: (req.user as User).id,
    name: (req.user as User).username,
    email: (req.user as User).email,
    avatar: (req.user as User).avatar,
    avatarUrl: `https://cdn.discordapp.com/avatars/${
      (req.user as User).discordId
    }/${(req.user as User).avatar}.png`,

    subscription: {
      status: subscription?.status ?? "none",
      plan: subscription?.type ?? "none",
    },
  });
});

// eslint-disable-next-line @typescript-eslint/no-misused-promises -- This is a callback function
app.post("/subscribe", async (req, res) => {
  try {
    if (!req.user) {
      res.status(401).send("Unauthorized");
      return;
    }

    const user: User = req.user as User;
    const { plan }: { plan: string } = req.body;

    const paymentSession = await stripe.checkout.sessions.create({
      payment_method_types: ["card"],
      line_items: [
        {
          price_data: {
            currency: "usd",
            product_data: {
              name: `Subscription - ${plan}`,
            },
            unit_amount: 500, // $5
            recurring: {
              interval: "month",
            },
          },
          quantity: 1,
        },
      ],
      mode: "subscription",
      success_url:
        "http://localhost:3000/success?session_id={CHECKOUT_SESSION_ID}",
      cancel_url: "http://localhost:3000/cancel",
    });

    await createSubscription({
      userId: user.id,
      plan,
      stripeSessionId: paymentSession.id,
      price: 500,
      status: "pending",
    });

    res.status(200).json({ sessionId: paymentSession.id });
  } catch (error) {
    res.status(500).send("Internal server error");
  }
});
app.use("/webhook", stripeWebhook);

app.use("/scan", scanRouter);

app.listen(3000, () => {
  // eslint-disable-next-line no-console -- This is a server
  console.log("Server is running on port 3000");
});
