<script>
const defineDescriptor = (src, dest, name) => {
  // eslint-disable-next-line no-prototype-builtins
  if (!dest.hasOwnProperty(name)) {
    const descriptor = Object.getOwnPropertyDescriptor(src, name);

    Object.defineProperty(dest, name, descriptor);
  }
};

const merge = (objs) => {
  const res = {};

  // eslint-disable-next-line guard-for-in
  for (const obj of objs) {
    const objKeys = obj
      ? Object.getOwnPropertyNames(obj)
      : [];

    // eslint-disable-next-line guard-for-in
    for (const objProp of objKeys) {
      defineDescriptor(obj, res, objProp);
    }
  }

  return res;
};

const buildFromProps = (obj, props) => {
  const res = {};

  // eslint-disable-next-line guard-for-in
  for (const prop of props) {
    defineDescriptor(obj, res, prop);
  }

  return res;
};

export default {
  props: {
    template: {
      type: String,
      required: false,
    },
    parent: {
      type: Object,
      required: false,
    },
    templateProps: {
      type: Object,
      default: () => ({}),
    },
  },
  render(h) {
    if (!this.template) {
      return null;
    }

    const parent = this.parent || this.$parent;

    const {
      $props: parentProps = {},
      $data: parentData = {},
      $options: parentOptions = {},
    } = parent;
    const {
      methods: parentMethods = {},
    } = parentOptions;
    const {
      $data = {},
      $props = {},
      $options: { methods = {} } = {},
    } = this;

    // build new objects by removing keys if already exists (e.g. created by mixins)
    const dataFromParent = {};
    for (const name in parentData) {
      if (typeof $data[name] === 'undefined') {
        dataFromParent[name] = parentData[name];
      }
    }

    const propsFromParent = {};
    for (const name in parentProps) {
      if (typeof $props[name] === 'undefined') {
        propsFromParent[name] = parentProps[name];
      }
    }

    const methodsFromParent = {};
    for (const name in parentMethods) {
      if (typeof methods[name] === 'undefined') {
        methodsFromParent[name] = parentMethods[name];
      }
    }

    const dataKeys = Object.keys(dataFromParent);
    const propsKeys = Object.keys(parentOptions.props ?? {});
    const methodKeys = Object.keys(methodsFromParent);
    const templatePropsKeys = Object.keys(this.templateProps);

    const propsTypes = [...dataKeys, ...propsKeys, ...methodKeys, ...templatePropsKeys];

    const methodsFromProps = buildFromProps(parent, methodKeys);

    const finalProps = merge([
      dataFromParent,
      propsFromParent,
      methodsFromProps,
      this.templateProps,
    ]);

    // eslint-disable-next-line no-underscore-dangle
    const provide = this.$parent._provided;

    const dynamic = {
      template: this.template || '<div></div>',
      props: propsTypes,
      computed: parentOptions.computed,
      components: parentOptions.components,
      provide,
    };

    // eslint-disable-next-line consistent-return
    return h(dynamic, {
      props: finalProps,
      on: this.$listeners,
    });
  },
};
</script>
