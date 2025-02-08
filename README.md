# creative coding bookclub

The home of the [creativecodingbook.club](https://creativecodingbook.club) website

# set up local development 

- [clone this repo to your machine](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository)
- Install the latest LTS version of [NodeJS](https://nodejs.org/en/learn/getting-started/introduction-to-nodejs) in your computer, ideally `v.22` or greater.
- Install [VSCode](https://code.visualstudio.com/) or use any IDE of your choice 
- Open a [terminal window](https://www.freecodecamp.org/news/command-line-for-beginners/) on the directory you cloned the repository, and run `npm install` to install all dependencies. This will create a `node_modules` folder. 
- Once the installation is done, on the same window, run `npm run dev`. This will start a development server running on `localhost:4321/`
- Open a tab in your browser and navigate to `localhost:4321/`
- Change any code in your IDE and see the changes reflected in real time on your browser!
- Once you're done `CTRL+C` inside your terminal window should stop the running process.

# upload your work 

If you're part of the bookclub and want to upload your work (or any other HTML pages) to [creativecodingbook.club](https://creativecodingbook.club) here is how:

We're using [Astro](https://docs.astro.build/en/concepts/why-astro/) to build the site, but its [template syntax](https://docs.astro.build/en/reference/astro-syntax/) should be very easy to grasp even if you have never played with it.

Start by creating a folder under `src/members/{your-handle}`. 

Every `.astro` file you create in that directory will be automatically routed to have its own page under [/{your-handle}/{page-name}](https://creativecodingbook.club/example/p5js-cdn) - see [the source code of this page](https://github.com/sb-luis/creative-coding-bookclub/tree/main/src/members/example) as an example.

To have your `.astro` page properly routed there is just one more step for you to do: you need to add a `.json` file [with the same filename as your page](https://github.com/sb-luis/creative-coding-bookclub/blob/main/src/members/sb-luis/hola.json) to be used as [page metadata](https://creativecodingbook.club/sb-luis/hola). This way you can set at least the `title`, `description`, and `keywords` of your page to help with SEO. We'll also use this metadata in the future to organize the pages per book and chapter.

Any `.js` files in your folder will be ignored, but you can use it to organize your code.

To keep things simple try only uploading `.astro`, `.js` or `.json` files for now. No images, or other assets. We can have a think later about how to organise these if we need them.

That's all! Go ahead and [fork this repository](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo), make your own folder and start [creating PRs](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork) to upload your creative coding sketches and custom pages to our site! 

If you post any bookclub-related content in social media feel free to mention us or tag the[#CreativeCodingBookclub](https://bsky.app/hashtag/CreativeCodingBookclub) hashtag for us to see your posts! 

If you get stuck or have any questions, join [C3S's Discord server](https://discord.gg/ggYbapqx) and message us at the `#creative-coding-bookclub` channel in it. 

# resources 

- [Freecodecamp - CLIs for begginers](https://www.freecodecamp.org/news/command-line-for-beginners/)
- [Astro - Syntax](https://docs.astro.build/en/reference/astro-syntax/)
- [Github - Forking a repository](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo)
- [Github - Collaborating with pull requests](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork)
- [Github - Git for begginers](https://github.blog/developer-skills/programming-languages-and-frameworks/what-is-git-our-beginners-guide-to-version-control/)