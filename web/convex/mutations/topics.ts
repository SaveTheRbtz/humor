import { mutation } from "../_generated/server";
import { v } from "convex/values";

export const upsertTopic = mutation({
  args: { name: v.string() },
  handler: async (ctx, args) => {
    const existingTopic = await ctx.db.query("topic").filter(q => q.eq(q.field("name"), args.name)).first();
    if (existingTopic) {
      return existingTopic._id;
    }
    return await ctx.db.insert("topic", { name: args.name });
  }
});