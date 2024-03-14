import { GRID_SIZES } from '@/constants';

export const getFlexPropsForGridRangeSize = (gridRangeSize) => {
  const [start, end] = gridRangeSize ?? [GRID_SIZES.min, GRID_SIZES.max];

  return [
    `offset-xs${start}`,
    `xs${end - start}`,
  ];
};
