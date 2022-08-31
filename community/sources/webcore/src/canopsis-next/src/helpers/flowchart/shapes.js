import { LINE_TYPES, SHAPES } from '@/constants';

import uid from '@/helpers/uid';

import { generatePoint } from './points';

/**
 * @typedef {Object} DefaultShapeConnection
 * @property {string} shapeId
 * @property {string} pointId
 * @property {string} side
 */

/**
 * @typedef {Object} DefaultShape
 * @property {string} _id
 * @property {string} type
 * @property {DefaultShapeConnection[]} connections
 * @property {string[]} connectedTo
 * @property {string} text
 * @property {Object} textProperties
 * @property {Object} properties
 * @property {boolean} aspectRatio
 */

/**
 * @typedef {DefaultShape} RectShape
 * @property {number} width
 * @property {number} height
 * @property {number} x
 * @property {number} y
 */

/**
 * @typedef {DefaultShape} LineShape
 * @property {Point[]} points
 * @property {string} lineType
 */

/**
 * @typedef {LineShape} ArrowLineShape
 */

/**
 * @typedef {LineShape} BidirectionalArrowLineShape
 */

/**
 * @typedef {DefaultShape} CircleShape
 * @property {number} diameter
 * @property {number} x
 * @property {number} y
 */

/**
 * @typedef {RectShape} EllipseShape
 */

/**
 * @typedef {RectShape} RhombusShape
 */

/**
 * @typedef {RectShape} ParallelogramShape
 * @property {number} offset
 */

/**
 * @typedef {RectShape} ProcessShape
 * @property {number} offset
 */

/**
 * @typedef {RectShape} DocumentShape
 * @property {number} offset
 */

/**
 * @typedef {RectShape} StorageShape
 * @property {number} radius
 */

/**
 * @typedef {RectShape} ImageShape
 * @property {number} src
 */

/**
 * @typedef {
 *   RectShape |
 *   CircleShape |
 *   RhombusShape |
 *   EllipseShape |
 *   ParallelogramShape |
 *   ImageShape |
 *   StorageShape |
 *   DocumentShape |
 *   ProcessShape |
 *   LineShape |
 *   ArrowLineShape |
 *   BidirectionalArrowLineShape
 * } Shape
 */

/**
 * Generate default shape
 *
 * @param {DefaultShape} shape
 * @returns {DefaultShape}
 */
const generateDefaultShape = shape => ({
  _id: shape._id ?? uid(),
  connections: shape.connections ?? [],
  connectedTo: shape.connectedTo ?? [],
  text: shape.text ?? '',
  textProperties: {
    ...shape.textProperties,
    fontColor: 'black',
    fontSize: 12,
  },
  properties: shape.properties ?? {},
  aspectRatio: shape.aspectRatio ?? false,
});

/**
 * Generate default rectangle shape
 *
 * @param {RectShape} shape
 * @returns {RectShape}
 */
export const generateRectShape = shape => ({
  ...generateDefaultShape(shape),

  type: SHAPES.rect,
  width: shape.width ?? 130,
  height: shape.height ?? 130,
  x: shape.x ?? 0,
  y: shape.y ?? 0,
});

/**
 * Generate default line shape
 *
 * @param {LineShape} shape
 * @returns {LineShape}
 */
export const generateLineShape = shape => ({
  ...generateDefaultShape(shape),

  type: SHAPES.line,
  lineType: shape.lineType ?? LINE_TYPES.line,
  points: shape.points ?? [
    generatePoint({
      x: 0,
      y: 100,
    }),
    generatePoint({
      x: 100,
      y: 0,
    }),
  ],
});

/**
 * Generate default arrow line shape
 *
 * @param {ArrowLineShape} shape
 * @returns {ArrowLineShape}
 */
export const generateArrowLineShape = shape => ({
  ...generateLineShape(shape),

  type: SHAPES.arrowLine,
});

/**
 * Generate default bidirectional arrow line shape
 *
 * @param {BidirectionalArrowLineShape} shape
 * @returns {BidirectionalArrowLineShape}
 */
export const generateBidirectionalArrowLineShape = shape => ({
  ...generateLineShape(shape),

  type: SHAPES.bidirectionalArrowLine,
});

/**
 * Generate default circle shape
 *
 * @param {CircleShape} shape
 * @returns {CircleShape}
 */
export const generateCircleShape = shape => ({
  ...generateDefaultShape(shape),

  type: SHAPES.circle,
  x: shape.x ?? 50,
  y: shape.y ?? 50,
  diameter: shape.diameter ?? 100,
});

/**
 * Generate default ellipse shape
 *
 * @param {EllipseShape} shape
 * @returns {EllipseShape}
 */
export const generateEllipseShape = shape => ({
  ...generateRectShape(shape),

  type: SHAPES.ellipse,
});

/**
 * Generate default rhombus shape
 *
 * @param {RhombusShape} shape
 * @returns {RhombusShape}
 */
export const generateRhombusShape = shape => ({
  ...generateRectShape(shape),

  type: SHAPES.rhombus,
});

/**
 * Generate default parallelogram shape
 *
 * @param {ParallelogramShape} shape
 * @returns {ParallelogramShape}
 */
export const generateParallelogramShape = shape => ({
  ...generateRectShape(shape),

  type: SHAPES.parallelogram,
  offset: shape.offset ?? 50,
});

/**
 * Generate default process shape
 *
 * @param {ProcessShape} shape
 * @returns {ProcessShape}
 */
export const generateProcessShape = shape => ({
  ...generateRectShape(shape),

  type: SHAPES.process,
  offset: shape.offset ?? 20,
});

/**
 * Generate default document shape
 *
 * @param {DocumentShape} shape
 * @returns {DocumentShape}
 */
export const generateDocumentShape = shape => ({
  ...generateRectShape(shape),

  type: SHAPES.document,
  offset: shape.offset ?? 20,
});

/**
 * Generate default storage shape
 *
 * @param {StorageShape} shape
 * @returns {StorageShape}
 */
export const generateStorageShape = shape => ({
  ...generateRectShape(shape),

  type: SHAPES.storage,
  radius: shape.radius ?? 20,
});

/**
 * Generate default image shape
 *
 * @param {ImageShape} shape
 * @returns {ImageShape}
 */
export const generateImageShape = shape => ({
  ...generateRectShape(shape),

  type: SHAPES.image,
  src: shape.src,
});
