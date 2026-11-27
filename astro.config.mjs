// @ts-check
import { defineConfig } from "astro/config";
import { satteri } from "@astrojs/markdown-satteri";
import starlight from "@astrojs/starlight";
import { cast } from "./src/lib/cast.mjs";

// SITE_URL is the address that the build is published at.
const site = process.env.SITE_URL ?? "https://usql.app";

export default defineConfig({
  site,
  trailingSlash: "ignore",
  output: "static",
  build: { format: "directory" },
  markdown: { processor: satteri({ mdastPlugins: [cast()] }) },
  integrations: [
    starlight({
      title: "usql",
      description: "A universal command line interface for SQL databases.",
      logo: { src: "./src/assets/logo.svg" },
      favicon: "/favicon.svg",
      social: [{ icon: "github", label: "GitHub", href: "https://github.com/xo/usql" }],
      editLink: { baseUrl: "https://github.com/xo/usql-site/edit/main/" },
      customCss: ["./src/styles/site.css"],
      components: { Head: "./src/components/Head.astro" },
      sidebar: [
        { label: "Get started", items: [{ autogenerate: { directory: "start" } }] },
        { label: "Guides", items: [{ autogenerate: { directory: "guides" } }] },
        { label: "Charts", items: [{ autogenerate: { directory: "charts" } }] },
        { label: "Project", items: [{ autogenerate: { directory: "project" } }] },
      ],
    }),
  ],
});
