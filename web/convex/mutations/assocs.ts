import { mutation } from "../_generated/server";
import { v } from "convex/values";

export const upsertAssocv1 = mutation({
  args: { text: v.string(), topic: v.id("topic"), model: v.id("model") },
  handler: async (ctx, args) => {
    const existingAssoc = await ctx.db.query("assocv1").filter(q => q.eq(q.field("text"), args.text)).first();
    if (existingAssoc) {
      return existingAssoc._id;
    }
    return await ctx.db.insert("assocv1", { text: args.text, topic: args.topic, model: args.model });
  }
});

export const upsertAssocv2 = mutation({
  args: { text: v.string(), topic: v.id("topic"), model: v.id("model") },
  handler: async (ctx, args) => {
    const existingAssoc = await ctx.db.query("assocv2").filter(q => q.eq(q.field("text"), args.text)).first();
    if (existingAssoc) {
      return existingAssoc._id;
    }
    return await ctx.db.insert("assocv2", { text: args.text, topic: args.topic, model: args.model });
  }
});

export const upsertAssocv3 = mutation({
  args: { text: v.string(), topic: v.id("topic"), model: v.id("model") },
  handler: async (ctx, args) => {
    const existingAssoc = await ctx.db.query("assocv3").filter(q => q.eq(q.field("text"), args.text)).first();
    if (existingAssoc) {
      return existingAssoc._id;
    }
    return await ctx.db.insert("assocv3", { text: args.text, topic: args.topic, model: args.model });
  }
});