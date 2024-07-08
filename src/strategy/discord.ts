import DiscordStrategy from "passport-discord";
import passport from "passport";
import { config } from "dotenv";
import type { User } from "@prisma/client";
import {
  createUser,
  findUserByDiscordId,
  getUserById,
} from "~/controller/user.controller";

config();

const scopes = ["identify", "email"];

passport.use(
  new DiscordStrategy(
    {
      clientID: process.env.CLIENT_ID ?? "",
      clientSecret: process.env.CLIENT_SECRET ?? "",
      callbackURL: process.env.REDIRECT_URI ?? "",
      scope: scopes,
    },
    // eslint-disable-next-line @typescript-eslint/no-misused-promises -- This is a callback function
    async (
      _accessToken,
      _refreshToken,
      profile: DiscordStrategy.Profile & { refreshToken: string },
      done
    ) => {
      profile.refreshToken = _refreshToken;
      // Check if the profile exists.
      const user = await findUserByDiscordId(profile.id);
      if (user) {
        done(null, user);
        return;
      }

      // Create a new user.
      const newUser = await createUser({
        discordId: profile.id,
        username: profile.username,
        email: profile.email ?? "",
        accessToken: _accessToken,
        avatar: profile.avatar ?? "",
        createdAt: new Date(),
        provider: profile.provider,
        updatedAt: new Date(),
      });
      done(null, newUser);
    }
  )
);

passport.serializeUser((user: User, done) => {
  done(null, user.id);
});

// eslint-disable-next-line @typescript-eslint/no-misused-promises -- This is a callback function
passport.deserializeUser(async (id: number, done) => {
  try {
    const user = await getUserById(id);
    done(null, user);
  } catch (error) {
    done(error, null);
  }
});
