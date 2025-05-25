/** @jsxImportSource @compiled/react */

interface RedCrossProps {
    left: number;
}

const CROSS_SIZE = 20;

export const RedCross = ({ left }: RedCrossProps) => (
    <svg
        css={{
            position: 'absolute',
            pointerEvents: 'none',
            zIndex: 999999,
            shapeRendering: 'geometricPrecision',
        }}
        style={{
            left: `${left}px`,
        }}
        width={CROSS_SIZE}
        height={CROSS_SIZE}
    >
        <line x1="0" y1="0" x2={CROSS_SIZE} y2={CROSS_SIZE} stroke="red" strokeWidth="3" />
        <line x1={CROSS_SIZE} y1="0" x2="0" y2={CROSS_SIZE} stroke="red" strokeWidth="3" />
    </svg>
);
