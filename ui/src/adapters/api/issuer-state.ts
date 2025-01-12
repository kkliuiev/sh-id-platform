import axios from "axios";
import dayjs from "dayjs";
import { z } from "zod";

import { APIResponse, ResultOK, buildAPIError, buildAuthorizationHeader } from "src/adapters/api";
import { getListParser, getStrictParser } from "src/adapters/parsers";
import { Env, IssuerStatus, Transaction, TransactionStatus } from "src/domain";
import { API_VERSION } from "src/utils/constants";
import { List } from "src/utils/types";

const transactionStatusParser = getStrictParser<TransactionStatus>()(
  z.union([
    z.literal("created"),
    z.literal("failed"),
    z.literal("pending"),
    z.literal("published"),
    z.literal("transacted"),
  ])
);

const transactionParser = getStrictParser<Transaction>()(
  z.object({
    id: z.number(),
    publishDate: z.coerce.date(z.string().datetime()),
    state: z.string(),
    status: transactionStatusParser,
    txID: z.string(),
  })
);

export async function publishState({ env }: { env: Env }): Promise<APIResponse<boolean>> {
  try {
    await axios({
      baseURL: env.api.url,
      headers: {
        Authorization: buildAuthorizationHeader(env),
      },
      method: "POST",
      url: `${API_VERSION}/state/publish`,
    });

    return { data: true, isSuccessful: true };
  } catch (error) {
    return { error: buildAPIError(error), isSuccessful: false };
  }
}

export async function retryPublishState({ env }: { env: Env }): Promise<APIResponse<boolean>> {
  try {
    await axios({
      baseURL: env.api.url,
      headers: {
        Authorization: buildAuthorizationHeader(env),
      },
      method: "POST",
      url: `${API_VERSION}/state/retry`,
    });

    return { data: true, isSuccessful: true };
  } catch (error) {
    return { error: buildAPIError(error), isSuccessful: false };
  }
}

export async function getStatus({
  env,
  signal,
}: {
  env: Env;
  signal?: AbortSignal;
}): Promise<APIResponse<IssuerStatus>> {
  try {
    const response = await axios({
      baseURL: env.api.url,
      headers: {
        Authorization: buildAuthorizationHeader(env),
      },
      method: "GET",
      signal,
      url: `${API_VERSION}/state/status`,
    });
    const { data } = resultOKStatusParser.parse(response);

    return { data, isSuccessful: true };
  } catch (error) {
    return { error: buildAPIError(error), isSuccessful: false };
  }
}

const resultOKStatusParser = getStrictParser<ResultOK<IssuerStatus>>()(
  z.object({
    data: z.object({ pendingActions: z.boolean() }),
    status: z.literal(200),
  })
);

export async function getTransactions({
  env,
  signal,
}: {
  env: Env;
  signal?: AbortSignal;
}): Promise<APIResponse<List<Transaction>>> {
  try {
    const response = await axios({
      baseURL: env.api.url,
      headers: {
        Authorization: buildAuthorizationHeader(env),
      },
      method: "GET",
      signal,
      url: `${API_VERSION}/state/transactions`,
    });
    const { data } = resultOKTransactionListParser.parse(response);

    return {
      data: {
        failed: data.failed,
        successful: data.successful.sort(
          ({ publishDate: a }, { publishDate: b }) => dayjs(b).unix() - dayjs(a).unix()
        ),
      },
      isSuccessful: true,
    };
  } catch (error) {
    return { error: buildAPIError(error), isSuccessful: false };
  }
}

const resultOKTransactionListParser = getStrictParser<
  ResultOK<unknown[]>,
  ResultOK<List<Transaction>>
>()(
  z.object({
    data: getListParser(transactionParser),
    status: z.literal(200),
  })
);
