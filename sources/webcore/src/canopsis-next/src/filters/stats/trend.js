export function trendColor(trendValue) {
  if (trendValue > 0) {
    return 'primary';
  } else if (trendValue < 0) {
    return 'error';
  }

  return 'black';
}

export function trendIcon(trendValue) {
  if (trendValue > 0) {
    return 'trending_up';
  } else if (trendValue < 0) {
    return 'trending_down';
  }

  return 'trending_flat';
}

