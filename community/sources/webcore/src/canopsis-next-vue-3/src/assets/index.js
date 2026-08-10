const modules = import.meta.glob('@/assets/event-filter-actions-types-images/*', {
  eager: true,
  import: 'default',
});

export const eventFilterActionsTypesImages = Object.fromEntries(
  Object.entries(modules).map(([path, url]) => {
    const name = path
      .split('/')
      .pop()
      .replace(/\.\w+$/, '');

    return [name, url];
  }),
);
