import { mutation } from "../_generated/server";
import { v } from "convex/values";

export const upsertPolicy = mutation({
  args: { text: v.string(), model: v.id("model"), name: v.string()},
  handler: async (ctx, args) => {
    const existingPolicy = await ctx.db.query("humor_policy").filter(q => q.eq(q.field("text"), args.text)).first();
    if (existingPolicy) {
      return existingPolicy._id;
    }
    return await ctx.db.insert("humor_policy", {
        name: args.name,
        text: args.text, 
        model: args.model,
    });
  }
});