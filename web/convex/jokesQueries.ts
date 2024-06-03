import { query } from "./_generated/server";
import { v } from "convex/values";

export const getRandomJokes = query({
  args: { count: v.number() },
  handler: async (ctx, args) => {
    const jokes = await ctx.db.query("jokes").collect();
    const shuffled = jokes.sort(() => 0.5 - Math.random());
    return shuffled.slice(0, args.count);
  },
});

export const getTopJokes = query({
  args: {},
  handler: async (ctx) => {
    const jokes = await ctx.db.query("jokes").collect();
    const filteredJokes = jokes.filter(joke => joke.views > 0n); // Ensure views are greater than zero
    const sortedJokes = filteredJokes.sort((a, b) => {
      const aScorePerView = Number(a.score) / Number(a.views);
      const bScorePerView = Number(b.score) / Number(b.views);
      return bScorePerView - aScorePerView;
    }).slice(0, 100);

    // Fetch topics and policies
    const topicIds = [...new Set(sortedJokes.map(joke => joke.topic))];
    const policyIds = [...new Set(sortedJokes.map(joke => joke.policy))];
    const assocv1 = [...new Set(sortedJokes.map(joke => joke.assocv1))];
    const assocv2 = [...new Set(sortedJokes.map(joke => joke.assocv2))];
    const assocv3 = [...new Set(sortedJokes.map(joke => joke.assocv3))];

    const topics = await Promise.all(topicIds.map(id => ctx.db.get(id)));
    const policies = await Promise.all(policyIds.map(id => ctx.db.get(id)));
    const assocv1s = await Promise.all(assocv1.map(id => ctx.db.get(id)));
    const assocv2s = await Promise.all(assocv2.map(id => ctx.db.get(id)));
    const assocv3s = await Promise.all(assocv3.map(id => ctx.db.get(id)));


    const validTopics = topics.filter(topic => topic !== null);
    const validPolicies = policies.filter(policy => policy !== null);
    const validAssocv1s = assocv1s.filter(assocv1 => assocv1 !== null);
    const validAssocv2s = assocv2s.filter(assocv2 => assocv2 !== null);
    const validAssocv3s = assocv3s.filter(assocv3 => assocv3 !== null);

    const topicMap = new Map(validTopics.map(topic => [topic?._id, topic?.name]));
    const policyMap = new Map(validPolicies.map(policy => [policy?._id, { name: policy?.name, text: policy?.text }]));
    const assocv1Map = new Map(validAssocv1s.map(assocv1 => [assocv1?._id, assocv1?.text]));
    const assocv2Map = new Map(validAssocv2s.map(assocv2 => [assocv2?._id, assocv2?.text]));
    const assocv3Map = new Map(validAssocv3s.map(assocv3 => [assocv3?._id, assocv3?.text]));


    
    // Attach topic names and policy info to jokes
    const jokesWithDetails = sortedJokes.map(joke => ({
      ...joke,
      topicName: topicMap.get(joke.topic),
      policyName: policyMap.get(joke.policy)?.name,
      policyText: policyMap.get(joke.policy)?.text,
      assocv1Text: assocv1Map.get(joke.assocv1),
      assocv2Text: assocv2Map.get(joke.assocv2),
      assocv3Text: assocv3Map.get(joke.assocv3),
    }));

    return jokesWithDetails;
  },
});