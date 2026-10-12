<template>
  <div ref="controlElement">
    <slot />
  </div>
</template>

<script>
import { DomEvent, control } from 'leaflet';

export default {
  inject: ['geomapRoot'],

  props: {
    position: {
      type: String,
      default: 'topright',
    },
  },

  data() {
    return {
      mapObject: undefined,
    };
  },

  mounted() {
    const element = this.$refs.controlElement;

    DomEvent.disableClickPropagation(element);
    DomEvent.disableScrollPropagation(element);

    const CustomControl = control({ position: this.position });

    CustomControl.onAdd = () => element;
    CustomControl.addTo(this.geomapRoot.mapObject);
    this.mapObject = CustomControl;
  },

  beforeUnmount() {
    this.mapObject?.remove();
  },
};
</script>
