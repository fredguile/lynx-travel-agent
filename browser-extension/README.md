# browser-extension

Firefox/Chrome extension to use as Lynx Travel Agent.

## Prerequisites

```sh
yarn install
```

## How to build

```sh
yarn build
```

The extension is packaged into the `dist` folder.

## How to develop

```sh
yarn watch
```

As Firefox user, I usually go to `about:debugging` and load the `dist` folder as temporary extension. 

When using `yarn watch`, Parcel will automatically reload the current page to force reload the latest extension code (HMR).
