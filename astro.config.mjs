// @ts-check
import { defineConfig } from "astro/config";
import { satteri } from "@astrojs/markdown-satteri";
import starlight from "@astrojs/starlight";
import starlightLinksValidator from "starlight-links-validator";
import { cast } from "./src/lib/cast.mjs";
import { releases } from "./src/lib/releases.mjs";
import { schemes } from "./src/lib/schemes.mjs";
import { tabs } from "./src/lib/tabs.mjs";

// SITE_URL is the address that the build is published at.
const site = process.env.SITE_URL ?? "https://usql.app";

export default defineConfig({
  site,
  trailingSlash: "ignore",
  output: "static",
  build: { format: "directory" },
  markdown: { processor: satteri({ mdastPlugins: [cast(), schemes(), tabs(), releases()] }) },
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
      // A link to a page that does not exist stops the build. A link to a file
      // of public/ passes. The editor at /admin/ is a directory of public/
      // that scripts/prepare.mjs makes, so it is named here.
      plugins: [starlightLinksValidator({ exclude: ["/admin", "/admin/"] })],
      sidebar: [
        { label: "Get started", items: [{ autogenerate: { directory: "start" } }] },
        { label: "Guides", items: [{ autogenerate: { directory: "guides" } }] },
        { label: "Charts", items: [{ autogenerate: { directory: "charts" } }] },
        { label: "Databases", items: [{ autogenerate: { directory: "databases" } }] },
        { label: "Project", items: [{ autogenerate: { directory: "project" } }] },
      ],
    }),
  ],
});
