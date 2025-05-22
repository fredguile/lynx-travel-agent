export function createRedCrossElement(event: MouseEvent) {
    // Draw a red cross at the cursor position
    const crossSize = 20;
    const cross = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    cross.setAttribute('width', `${crossSize}`);
    cross.setAttribute('height', `${crossSize}`);
    cross.style.position = 'absolute';
    cross.style.left = `${event.pageX - crossSize / 2}px`;
    cross.style.top = `${event.pageY - crossSize / 2}px`;
    cross.style.pointerEvents = 'none';
    cross.style.zIndex = '999999';
    cross.innerHTML = `
        <line x1="0" y1="0" x2="${crossSize}" y2="${crossSize}" stroke="red" stroke-width="3" />
        <line x1="${crossSize}" y1="0" x2="0" y2="${crossSize}" stroke="red" stroke-width="3" />
    `;

    return cross;
}