export const mapInformationPopupMixin = {
  props: {
    iconSize: {
      type: Number,
      default: 24,
    },
    popupTemplate: {
      type: String,
      required: false,
    },
    popupActions: {
      type: Boolean,
      default: false,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    pbehaviorEnabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      positionX: 0,
      positionY: 0,
      activePoint: undefined,
    };
  },
  methods: {
    openPopup(point, event) {
      const { top, left, width } = event.target.getBoundingClientRect();

      this.positionY = top;
      this.positionX = left + width / 2;
      this.activePoint = point;
    },

    closePopup() {
      this.activePoint = undefined;
    },

    showLinkedMap() {
      this.$emit('show:map', this.activePoint.map);
      this.closePopup();
    },

    showAlarms() {
      this.$emit('show:alarms', this.activePoint);
      this.closePopup();
    },
  },
};
