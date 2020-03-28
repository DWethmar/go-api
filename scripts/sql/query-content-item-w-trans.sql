SELECT 
id, 
name, 
json_object_agg(t.locale, t.attrs) as attrs,
created_on, 
updated_on
FROM public.content_item c
LEFT JOIN public.content_item_translation t ON c.id = t.content_item_id AND t.locale IN ('nl', 'en')
WHERE c.id = '5e24490d-56c8-4a2d-a023-9f978f636bef'
GROUP BY c.id



SELECT 
id, 
name, 
created_on, 
updated_on,
COALESCE(jsonb_object_agg(t.locale, t.attrs) FILTER (WHERE t.locale IS NOT NULL), '{}'::JSONB) as attrs

FROM public.content_item c
LEFT JOIN public.content_item_translation t ON c.id = t.content_item_id AND t.locale IN ('nl', 'en')

GROUP BY c.id
