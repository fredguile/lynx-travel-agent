import { createLogger } from "../utils";

const ENDPOINT = `http://localhost:5678/webhook/1831ad0f-9c9b-4fb3-99e6-1ce8c0857931`;

const log = createLogger('analyseScreenContext');

export async function analyseScreenContext(currentUrl: string, screenshot: Blob) {
    log('analysing screen context', currentUrl);

    const formData = new FormData();
    formData.append('screenshot', screenshot, 'screenshot.png');
    let res = await fetch(`${ENDPOINT}?currentUrl=${encodeURIComponent(currentUrl)}`, {
        method: 'POST',
        body: formData,
    });
    const { fields } = await res.json();

    log('got screen context', { fields });

    return {
        fields: fields.map((field: any) => ({
            label: field.label,
            description: field.description,
            inputType: field.inputType,
            required: field.required,
        }))
    };
}
