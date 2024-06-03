import { mutation } from "../_generated/server";
import { v } from "convex/values";

export const upsertModel = mutation({
  args: { name: v.string(), temperature: v.number() },
  handler: async (ctx, args) => {
    const existingModel = await ctx.db.query("model").filter(q => q.eq(q.field("name"), args.name)).first();
    if (existingModel) {
      return existingModel._id;
    }
    return await ctx.db.insert("model", { name: args.name, temperature: args.temperature });
  }
});