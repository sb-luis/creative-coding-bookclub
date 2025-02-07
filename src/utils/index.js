export function getMemberPagesRoutes() {
  // Import `.astro` file member pages
  const memberPages = import.meta.glob('../members/**/*.astro', {
    eager: true,
  });
  let routes = Object.keys(memberPages);

  // Import `.json` metadata overrides 
  const memberPagesMetadata = import.meta.glob(`../members/**/**.json`, {
    eager: true,
  });

  // Merge default metadata with overrides
  routes = Object.keys(memberPages)?.map((key) => {
    const alias = key.split('/')[2];
    const page = key.split('/')[3].split('.')[0];
    const metadataOverrides =
      memberPagesMetadata[`../members/${alias}/${page}.json`] ?? {};

    const metadata = {
      title: 'member page',
      description: 'member page description',
      keywords: 'Creative Coding, Bookclub, Member',
      ...metadataOverrides,
    };

    delete metadata.default;

    if (!Object.keys(metadataOverrides).length) {
      metadata.isValidPage = true;
    }

    return {
      params: {
        alias,
        page,
      },
      props: { component: memberPages[key], metadata },
    };
  });

  // Filter out any pages that don't have a `.json` file
  routes = routes.filter((p) => !p.props.metadata.isValidPage);

  return routes;
}
