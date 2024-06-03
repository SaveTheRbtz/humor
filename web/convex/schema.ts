import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
    humor_policy: defineTable({
        name: v.string(),
        text: v.string(),
        model: v.id("model"),
    }),
    assocv1: defineTable({
        text: v.string(),
        topic: v.id("topic"),
        model: v.id("model"),
    }),
    assocv2: defineTable({
        text: v.string(),
        topic: v.id("topic"),
        model: v.id("model"),
    }),
    assocv3: defineTable({
        text: v.string(),
        topic: v.id("topic"),
        model: v.id("model"),
    }),
    model: defineTable({
        name: v.string(),
        temperature: v.number(),
    }),
    topic: defineTable({
        name: v.string(),
    }),
    jokes: defineTable({
        text: v.string(),
        views: v.int64(),
        score: v.int64(),

        topic: v.id("topic"),
        model: v.id("model"),
        policy: v.id("humor_policy"),
        assocv1: v.id("assocv1"),
        assocv2: v.id("assocv2"),
        assocv3: v.id("assocv3"),
    })
    .index("by_topic", ["topic"])
});
