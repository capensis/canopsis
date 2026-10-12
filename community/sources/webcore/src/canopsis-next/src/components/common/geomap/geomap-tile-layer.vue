<script>
import { tileLayer } from 'leaflet';

export default {
  inject: ['geomapRoot'],

  props: {
    url: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      required: false,
    },
    layerType: {
      type: String,
      default: 'base',
    },
    noWrap: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      mapObject: undefined,
    };
  },

  mounted() {
    this.mapObject = tileLayer(this.url, {
      noWrap: this.noWrap,
    });

    this.geomapRoot.addLayer(this);
  },

  beforeUnmount() {
    this.geomapRoot.removeLayer(this);
  },

  render() {
    return null;
  },
};
</script>
