/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { useQuery } from '@tanstack/react-query'
import { FileWarning } from 'lucide-react'
import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'

import { PublicLayout } from '@/components/layout'
import { RichContent } from '@/components/rich-content'
import { Card, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { isHttpUrl, isLikelyHtml } from '@/lib/content-format'

import type { LegalDocumentResponse } from './types'

// SipToken: the website addresses its legal pages per locale; `{lang}` in the
// configured URL becomes the website's code for the console language.
function localizeDocumentUrl(url: string, language: string): string {
  if (!url.includes('{lang}')) return url
  const lang = language.toLowerCase()
  let site = 'en'
  if (
    ['zh-tw', 'zh-hk', 'zh-mo'].includes(lang) ||
    lang.startsWith('zh-hant')
  ) {
    site = 'zh-hant'
  } else if (lang.startsWith('zh')) {
    site = 'zh'
  } else if (lang.startsWith('ru')) {
    site = 'ru'
  }
  return url.replaceAll('{lang}', site)
}

type LegalDocumentProps = {
  title: string
  queryKey: string
  fetchDocument: () => Promise<LegalDocumentResponse>
  emptyMessage: string
}

export function LegalDocument({
  title,
  queryKey,
  fetchDocument,
  emptyMessage,
}: LegalDocumentProps) {
  const { i18n } = useTranslation()
  const { data, isLoading } = useQuery({
    queryKey: [queryKey],
    queryFn: fetchDocument,
    staleTime: 10 * 60 * 1000,
  })

  const rawContent = data?.data?.trim() ?? ''
  const hasContent = rawContent.length > 0
  const isUrl = hasContent && isHttpUrl(rawContent)
  const contentIsHtml = hasContent && isLikelyHtml(rawContent)
  const success = data?.success ?? false
  // SipToken: open an external document directly, in the console's language
  // (SIPTOKEN.md)
  const externalUrl =
    success && isUrl ? localizeDocumentUrl(rawContent, i18n.language) : ''

  useEffect(() => {
    if (externalUrl) window.location.replace(externalUrl)
  }, [externalUrl])

  if (isLoading || externalUrl) {
    return (
      <PublicLayout>
        <div className='mx-auto flex max-w-4xl flex-col gap-4 py-12'>
          <Skeleton className='h-8 w-[45%]' />
          <Skeleton className='h-4 w-full' />
          <Skeleton className='h-4 w-[90%]' />
          <Skeleton className='h-4 w-[80%]' />
        </div>
      </PublicLayout>
    )
  }

  if (!success || !hasContent) {
    return (
      <PublicLayout>
        <div className='mx-auto max-w-2xl py-12'>
          <Card className='border-dashed'>
            <CardHeader className='flex flex-row items-center gap-4'>
              <div className='bg-muted rounded-lg p-2'>
                <FileWarning className='text-muted-foreground h-5 w-5' />
              </div>
              <div className='space-y-1'>
                <CardTitle className='text-lg font-semibold'>{title}</CardTitle>
                <p className='text-muted-foreground text-sm'>
                  {data?.message || emptyMessage}
                </p>
              </div>
            </CardHeader>
          </Card>
        </div>
      </PublicLayout>
    )
  }

  return (
    <PublicLayout showMainContainer={!contentIsHtml}>
      {contentIsHtml ? (
        <RichContent mode='html' htmlVariant='isolated' content={rawContent} />
      ) : (
        <div className='mx-auto max-w-4xl space-y-6 py-12'>
          <div className='space-y-2'>
            <h1 className='text-3xl font-semibold tracking-tight'>{title}</h1>
          </div>

          <RichContent
            mode='markdown'
            content={rawContent}
            className='prose-neutral dark:prose-invert max-w-none'
          />
        </div>
      )}
    </PublicLayout>
  )
}
