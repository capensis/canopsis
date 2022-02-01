export const actionsPanelItem = {
  props: ['method'],
  template: `
      <button class="actions-panel-item" @click="method" />
    `,
};
