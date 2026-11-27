// The pages of the site are the Markdown files of src/content/docs. Their
// frontmatter is Starlight's, and the CMS fields of cms/ follow it. When you
// add a field, add it here and in cms/ in the same change. See docs/CONTENT.md.
import { defineCollection } from "astro:content";
import { docsLoader } from "@astrojs/starlight/loaders";
import { docsSchema } from "@astrojs/starlight/schema";

export const collections = {
  docs: defineCollection({ loader: docsLoader(), schema: docsSchema() }),
};
