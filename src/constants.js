const IS_PROD = import.meta.env.PROD

export const REPO_PATH = 'sb-luis/creative-coding-bookclub'
export const DOMAIN = IS_PROD ? 'https://creativecodingbook.club' : 'http://localhost:4321'
