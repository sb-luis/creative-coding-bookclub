# creative coding bookclub

The home of the [creativecodingbook.club](https://creativecodingbook.club) website.

## set up local development

### software requirements

- [git](https://git-scm.com/downloads)
- the latest LTS version of [NodeJS](https://nodejs.org/en/learn/getting-started/introduction-to-nodejs), ideally `v.22` or greater.
- some terminal emulator, whatever is built-in in your computer should be fine
- some text or code editor, like [VSCode](https://code.visualstudio.com/)
  - _optional (but recommended with VSCode):_ [Astro's VSCode extension](https://marketplace.visualstudio.com/items?itemName=astro-build.astro-vscode)

### steps

- [fork this repository](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo) to your github account
- open a terminal window and navigate to the folder where you want to store the code
- [clone](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository) the newly forked repository in that folder
- run `npm install` to install all dependencies. This will create a `node_modules` folder
- once the installation is done, run `npm run dev`. This will start a development server running on `localhost:4321/`
- open a tab in your browser and navigate to `localhost:4321/`
- change any code and see the changes reflected in real time on your browser!
- once you're done `CTRL+C` inside your terminal window should stop the running NodeJS process

## uploading your work

If you're part of the bookclub and want to upload your work (or any other HTML pages) to [creativecodingbook.club](https://creativecodingbook.club) start by deciding on an member's alias that is [not in use yet](https://github.com/sb-luis/creative-coding-bookclub/tree/main/src/members)

All your sketches or any other pages will be linked under this alias, like [creativecodingbook.club/{alias}/{page}](https://creativecodingbook.club/davide/sketch-00-1/).

Once you know the name you want to use, you have two ways to upload your work: writing in our Discord channel or opening a Pull Request.

### through Discord

Paste a link to the source code of your sketch in the `#creative-coding-bookclub` channel of [C3S's Discord server](https://discord.gg/ggYbapqx)

Let us know explicitely in the message you're posting that you would like your sketch visible on the website, otherwise we'll assume you're simply sharing it internally!

Pick at least a title for the sketch, and the url you would like for it (e.g. `Sketch 1` and `/sketch-1`).

If you want you can also let us know a short description for it and some keywords (for search engines).

We'll add the sketch to the website for you (or upload it / delete it depending on which was your request)

### through a Pull Request

If you are comfortable with web technologies and git, this method is always prefered! 

Start by setting up the repository in your local development environment by following the steps listed above.

Once you've managed to get the repository running locally, make a folder under `src/members/{your-alias}`.

We're using [Astro](https://docs.astro.build/en/concepts/why-astro/) to build the site, but its [template syntax](https://docs.astro.build/en/reference/astro-syntax/) should be hopefully very easy to grasp even if you have never played with it.

Every `.astro` file you create in your member's directory will be automatically routed to have its own page under [/{your-alias}/{page-name}](https://creativecodingbook.club/_example/p5js-cdn) - check the [example pages](https://github.com/sb-luis/creative-coding-bookclub/tree/main/src/members/_example) to get a better idea.

To have your `.astro` page properly routed there is just one more step for you to do: you need to add a `.json` file [with the same filename as your page](https://github.com/sb-luis/creative-coding-bookclub/blob/main/src/members/_example/metadata.json) to be used as [page metadata](https://creativecodingbook.club/_example/metadata). This way you can set at least the `title`, `description`, and `keywords` of your page to help with SEO. We'll also use this metadata in the future to organize the pages per book and chapter.

Any `.js` files in your folder will be ignored, but you can use them to organize your code.

To keep things simple try only uploading `.astro`, `.js` or `.json` files for now. No images, or other assets. We can have a think later about how to organise these if we need them.

> A note on using P5 js

If you're using [P5JS](https://p5js.org/) in your sketches, you can run it either in [global mode](https://github.com/sb-luis/creative-coding-bookclub/blob/main/src/members/_example/p5js-global-mode.astro) or [instance mode](https://github.com/sb-luis/creative-coding-bookclub/blob/main/src/members/_example/p5js-instance-mode.astro). Feel free to read [their docs](https://github.com/processing/p5.js/wiki/Global-and-instance-mode) for more info.

You might be familiar with the former if you have been using the [P5 Web Editor](https://editor.p5js.org/). If you want to use `global mode`, remember to add the [following line](https://github.com/sb-luis/creative-coding-bookclub/blob/main/src/members/_example/p5js-global-mode.json#L6) to your `.json` file.

Once your sketches grow in complexity, using the library in `instance mode` will help you organize your code around [modules](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules).

# resources

- [Freecodecamp - CLIs for begginers](https://www.freecodecamp.org/news/command-line-for-beginners/)
- [Astro - Syntax](https://docs.astro.build/en/reference/astro-syntax/)
- [Github - Forking a repository](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo)
- [Github - Collaborating with pull requests](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork)
- [Github - Git for begginers](https://github.blog/developer-skills/programming-languages-and-frameworks/what-is-git-our-beginners-guide-to-version-control/)
