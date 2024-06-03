import { mutation } from "./_generated/server";
import { v } from "convex/values";

export const updateJokes = mutation({
  args: { 
    jokeId1: v.id("jokes"), 
    jokeId2: v.id("jokes"), 
    winningJokeId: v.union(v.id("jokes"), v.null())
  },
  handler: async (ctx, args) => {
    const joke1 = await ctx.db.get(args.jokeId1);
    const joke2 = await ctx.db.get(args.jokeId2);
    
    if (!joke1 || !joke2) throw new Error("Joke not found");

    // Increment views for both jokes
    await ctx.db.patch(args.jokeId1, {
      views: joke1.views + BigInt(1),
    });
    await ctx.db.patch(args.jokeId2, {
      views: joke2.views + BigInt(1),
    });

    // Update scores based on the winning joke
    if (args.winningJokeId) {
      if (args.winningJokeId === args.jokeId1) {
        await ctx.db.patch(args.jokeId1, {
          score: joke1.score + BigInt(1),
        });
        await ctx.db.patch(args.jokeId2, {
          score: joke2.score - BigInt(1),
        });
      } else {
        await ctx.db.patch(args.jokeId1, {
          score: joke1.score - BigInt(1),
        });
        await ctx.db.patch(args.jokeId2, {
          score: joke2.score + BigInt(1),
        });
      }
    } else {
      await ctx.db.patch(args.jokeId1, {
        score: joke1.score - BigInt(1),
      });
      await ctx.db.patch(args.jokeId2, {
        score: joke2.score - BigInt(1),
      });
    }
  },
});