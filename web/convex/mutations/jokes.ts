import { mutation } from "../_generated/server";
import { v } from "convex/values";

export const insertJoke = mutation({
  args: {
    text: v.string(),
    views: v.int64(),
    score: v.int64(),
    model: v.id("model"),
    policy: v.id("humor_policy"),
    assocv1: v.id("assocv1"),
    assocv2: v.id("assocv2"),
    assocv3: v.id("assocv3"),
    topic: v.id("topic")
  },
  handler: async (ctx, args) => {
    return await ctx.db.insert("jokes", args);
  }
});