export const createActivatorElementStub = className => ({
  template: `
    <div class='${className}'>
      <slot name="activator" />
      <slot />
    </div>
  `,
});
