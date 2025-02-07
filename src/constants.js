const IS_PROD = import.meta.env.PROD;

export const DOMAIN = IS_PROD ? 'https://creativecodingbook.club' : 'http://localhost:4321'