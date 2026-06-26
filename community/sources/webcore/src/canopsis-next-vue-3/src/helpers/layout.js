/**
 * Inline flex styles matching the former v-flex-compat xsN column widths.
 *
 * @param {number} cols Grid columns out of 12.
 * @returns {{ flex: string, maxWidth: string, minWidth: number }}
 */
export const gridFlexStyle = cols => ({
  flex: `0 0 calc(100% * ${cols} / 12)`,
  maxWidth: `calc(100% * ${cols} / 12)`,
  minWidth: 0,
});
