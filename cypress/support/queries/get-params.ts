import { logCommand, updateConsoleProps } from "support/commands/helpers";

type GetParamsResult = Record<string, string | string[]>;

export default function (this: unknown, ...names: string[]) {
  const log = logCommand("getParams", getInitialConsoleProps(names), names);

  const urlFn = cy.now("url", { log: false }) as () => string;

  return () => {
    const result = getParamsResult(urlFn(), names);

    if (cy.state("current") === this) {
      updateConsoleProps(log, (cp) => {
        cp.yielded = result;
      });
    }

    return result;
  };
}

function getInitialConsoleProps(names: string[]): Cypress.ObjectLike {
  switch (names.length) {
    case 0:
      return {};

    case 1: {
      return { name: names.at(0) };
    }

    default:
      return { names };
  }
}

function getParamsResult(url: string, names: string[]) {
  const params = getParamsObject(url);

  switch (names.length) {
    case 0:
      return params;

    case 1:
      return params[names.at(0)!];

    default:
      return Cypress._.pick(params, ...names);
  }
}

function getParamsObject(url: string): GetParamsResult {
  const { searchParams } = new URL(url);

  return [...searchParams].reduce(
    (previous: GetParamsResult, [name, value]) => {
      if (name in previous) {
        const previousValue = previous[name];
        if (Array.isArray(previousValue)) {
          previousValue.push(value);
        } else {
          previous[name] = [value];
        }
      } else {
        previous[name] = value;
      }

      return previous;
    },
    {},
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      getParams(...names: string[]): Chainable<GetParamsResult>;

      getParams(name: string): Chainable<string | string[]>;

      getParams(): Chainable<GetParamsResult>;
    }
  }
}
