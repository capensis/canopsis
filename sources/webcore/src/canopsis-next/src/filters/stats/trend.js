export default function (trendValue) {
  if (trendValue > 0) {
    return {
      icon: 'trending_up',
      color: 'primary',
    };
  } else if (trendValue < 0) {
    return {
      icon: 'trending_down',
      color: 'error',
    };
  }

  return {
    icon: 'trending_flat',
  };
}
