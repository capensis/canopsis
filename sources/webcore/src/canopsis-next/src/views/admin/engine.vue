<template lang="pug">
  div
    h1 ENgine
    div
      v-stepper.engines-stepper(vertical)
        template(v-for="(engine, index) in engines")
          v-stepper-step.engine-step(
            :key="engine.name",
            :step="index"
          )
            div
              div
                strong FIFO
              small Manages the queue of events
          v-stepper-content(
            :key="`${engine.name}-content`",
            :step="index"
          )
</template>

<script>
export default {
  data() {
    return {
      engines: [
        {
          name: 'engine-webhook',
          input_queue: 'Engine_webhook',
          output_queue: 'Engine_action',
        },
        {
          name: 'engine-fifo',
          input_queue: 'Engine_fifo',
          output_queue: 'Engine_che',
        },
        {
          name: 'engine-axe',
          input_queue: 'Engine_axe',
          output_queue: 'Engine_watcher',
        },
        {
          name: 'engine-pbehavior',
          input_queue: 'Engine_pbehavior',
          output_queue: 'Engine_axe',
        },
        {
          name: 'engine-che',
          input_queue: 'Engine_che',
          output_queue: 'Engine_pbehavior',
        },
        {
          name: 'engine-action',
          input_queue: 'Engine_action',
          output_queue: '',
        },
        {
          name: 'engine-watcher',
          input_queue: 'Engine_watcher',
          output_queue: 'Engine_dynamic_infos',
        },
        {
          name: 'engine-dynamic-info',
          input_queue: 'Engine_dynamic_infos',
          output_queue: 'Engine_webhook',
        },
      ],
    };
  },
};
</script>

<style lang="scss" scoped>
.engines-stepper {
  & /deep/ .v-stepper__step__step {
    height: 40px;
    width: 40px;
    max-width: 40px;
    margin-right: 30px;
    font-size: 15px;
    font-weight: 500;
  }

  &.v-stepper--vertical /deep/ .v-stepper__content {
    margin: -8px -36px -16px 44px;
  }
}
</style>
